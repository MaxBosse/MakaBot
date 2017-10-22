package command

import (
	"sync"
	"time"

	"github.com/MaxBosse/MakaBot/cache"
	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
)

type Music struct {
	subCommands      map[string]Command
	MusicGuilds      map[string]chan *musicChan
	musicGuildsMutex sync.Mutex
	MusicRuntimes    map[string]*musicHandler
	parent           Command
}

type musicChan struct {
	guildID   string
	channelID string
	videoInfo *ytdl.VideoInfo
}

type musicHandler struct {
	Queue     []*ytdl.VideoInfo
	Done      chan error
	ChannelID string
	Playing   bool
	Stop      chan bool
	Skip      chan bool
}

func init() {
	cmd := new(Music)
	cmd.subCommands = make(map[string]Command)
	cmd.subCommands["play"] = new(MusicPlay)
	cmd.subCommands["stop"] = new(MusicStop)
	cmd.subCommands["skip"] = new(MusicSkip)
	cmd.subCommands["help"] = new(Help)

	cmd.MusicGuilds = make(map[string]chan *musicChan)
	cmd.MusicRuntimes = make(map[string]*musicHandler)

	Register(cmd)
}

func (t *Music) Name() string {
	return "music"
}

func (t *Music) Description() string {
	return "Allows playing music"
}

func (t *Music) Usage() string {
	return "[command]"
}

func (t *Music) SubCommands() map[string]Command {
	return t.subCommands
}

func (t *Music) Parent() Command {
	return t.parent
}

func (t *Music) SetParent(cmd Command) {
	t.parent = cmd
}

func (t *Music) Event(c *Context, event *discordgo.Event) {
	switch ty := event.Struct.(type) {
	case *discordgo.Ready:
		for _, g := range ty.Guilds {
			go t.initMusicQueue(c, g)
		}
	case *discordgo.GuildCreate:
		go t.initMusicQueue(c, ty.Guild)
	}
	return
}

func (t *Music) Message(c *Context) {
	log.Debugln(t.Name() + " called")
	if handleSubCommands(c, t) {
		return
	}

	serverConf, err := c.Cache.GetServer(c.Guild.ID)
	if err != nil {
		serverConf = cache.CacheServer{}
	}

	c.Send("Please use `" + serverConf.Prefix + t.Name() + " " + t.Usage() + "`\n" +
		"Use `" + serverConf.Prefix + "help " + t.Name() + "` for more info!")

}

func (t *Music) initMusicQueue(c *Context, g *discordgo.Guild) {
	t.musicGuildsMutex.Lock()
	if _, ok := t.MusicGuilds[g.ID]; ok {
		close(t.MusicGuilds[g.ID])
	}

	t.MusicGuilds[g.ID] = make(chan *musicChan, 5)
	t.musicGuildsMutex.Unlock()

	for mC := range t.MusicGuilds[g.ID] {
		if _, ok := t.MusicRuntimes[mC.channelID]; ok {
			t.MusicRuntimes[mC.channelID].Queue = append(t.MusicRuntimes[mC.channelID].Queue, mC.videoInfo)
			log.Noteln("Songs in queue:", len(t.MusicRuntimes[mC.channelID].Queue))
			continue
		}

		t.MusicRuntimes[mC.channelID] = &musicHandler{
			ChannelID: mC.channelID,
			Done:      make(chan error),
			Skip:      make(chan bool, 1),
			Stop:      make(chan bool, 1),
			Playing:   false,
		}
		t.MusicRuntimes[mC.channelID].Queue = append(t.MusicRuntimes[mC.channelID].Queue, mC.videoInfo)
		go t.musicQueue(c, g, t.MusicRuntimes[mC.channelID])
	}
}

func (t *Music) musicQueue(c *Context, g *discordgo.Guild, mH *musicHandler) {

	// Join the provided voice channel.
	vc, err := c.Session.ChannelVoiceJoin(g.ID, mH.ChannelID, false, true)
	if err != nil {
		log.Errorln(err)
	}

	timeout := time.NewTimer(time.Minute)
	for {

		select {
		case <-timeout.C:
			log.Noteln("Got timeout, disconnecting")
			vc.Disconnect()
			return
		default:
			if len(mH.Queue) > 0 {
				var currentSong *ytdl.VideoInfo
				currentSong, mH.Queue = mH.Queue[0], mH.Queue[1:]

				log.Noteln("Playing "+currentSong.Title, "Songs in queue:", len(mH.Queue))

				format := currentSong.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
				videoURL, err := currentSong.GetDownloadURL(format)
				if err != nil {
					log.Errorln(err)
					return
				}

				options := dca.StdEncodeOptions
				options.RawOutput = true
				options.Bitrate = 96
				options.Application = dca.AudioApplicationAudio
				options.Volume = 200
				options.Threads = 4

				encSesh, err := dca.EncodeFile(videoURL.String(), options)
				if err != nil {
					log.Errorln(err)
					return
				}
				defer encSesh.Cleanup()

				mH.Playing = true
				dca.NewStream(encSesh, vc, mH.Done)

				select {
				case err = <-mH.Done:
					log.Errorln(err)
					mH.Playing = false
					break
				case <-mH.Stop:
					log.Noteln("Stopped song")
					mH.Playing = false
					vc.Disconnect()
					return
				case <-mH.Skip:
					encSesh.Cleanup()
					<-mH.Done
					break
				}

				log.Noteln("Songs in queue:", len(mH.Queue))
				// Reset timeout
				timeout = time.NewTimer(time.Minute)
			} else {
				log.Noteln("No more songs in queue.")
				vc.Disconnect()
				return
			}
		}
	}
}
