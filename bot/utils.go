package bot

import (
	"github.com/bwmarrin/discordgo"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func RemoveDuplicateMembers(list *[]*discordgo.Member) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *list {
		if !found[x.User.ID] {
			found[x.User.ID] = true
			(*list)[j] = (*list)[i]
			j++
		}
	}
	*list = (*list)[:j]
}
