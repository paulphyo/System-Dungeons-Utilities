{{/*
	Recommended Trigger: Reaction Added
	Recommended Trigger Type: Reaction

	Usage: Execute auction start/stop events based on reactions

	Copyright (c): zen | ゼン#0008; 2022
	License: MIT
	Repository: https://github.com/z3nn13/System-Dungeons
*/}}

{{/* Global Variables */}}
{{$auction_channelID := 781498393112477726}}
{{$reminder_channelID := 794761162729848863}}
{{$treasurer_role := 860729740268011561}}

{{$keyWords := cslice "box" "key" "nitro" "helmet"}}
{{$pics := cslice "🎁" "<:keyY:860801927775125504>" "<:nitro:860798633368092673>"	"<:Helmet:710741625168461846>"}}

{{$data := sdict (or .ExecData sdict)}}
{{$auc_event := false}}
{{$bank_event := false}}

{{if not .ExecData}}
	{{if eq .Message.Author.ID .User.ID}}

		{{if eq ($k:=.Reaction.Emoji.Name) "⏰" "🗑️" }}
			{{with dbGet .Message.ID "can_auction" }}{{ dbDel .UserID .Key }}
				{{ if eq $k "⏰"}}
					{{ $auc_event = true }}
					{{ $data = .Value }}
					{{ deleteAllMessageReactions nil $.Message.ID }}
				{{ else }} {{/* 🗑️ */}}
					{{ deleteTrigger 1 }}
				{{ end }}
			{{ end }}
		
		{{ else if eq ($k:=.Reaction.Emoji.Name) "🏦" "🆎" }}
			{{ with dbGet .User.ID "ongoing" }}
				{{ $data = .Value }}
				{{ if eq $data.originalMsg (str $.Message.ID)}}
					{{ if eq $k "🆎" }}
						{{ dbDel .UserID .Key }}{{.UserID}} {{.Key}}
						{{ deleteAllMessageReactions nil $.Message.ID}}
						{{ $_ := exec "delreminder" $data.remID}}
						{{ $r := exec "remind" "1s" $data.remText }}
						{{ scheduleUniqueCC $.CCID nil 1 $.Message.ID $data}}
					{{ else }} {{/* 🏦 */}}
						{{ deleteAllMessageReactions nil $.Message.ID $k}}
						{{$bank_event = true}}
					{{ end}}
				{{ end}}
			{{end}}
		{{end}}

	{{else}}
		{{deleteMessageReaction nil .Message.ID .User.ID .Reaction.Emoji.APIName}}
	{{end}}
{{end}}


{{/* ⏰ */}}
{{if $auc_event}}
	{{ $time := $data.auction.Get "time period" | toDuration}}
	{{ $item := $data.auction.Get "item name + quantity"}}
	{{ $remID := ""}}{{$emoji := ""}}


	{{/* reminder stuff */}}
	{{ $beforeslice := split (exec "reminders") "\n"}}
	{{ $remText := joinStr "\n"
	"❄️ **Auction Over** ❄️" 
	(print "> Item: `" $item "`")
	(print "> Ending Time: <t:" (currentTime.Add $time).Unix ":R> <:stopwatch:845919763414777857>")
	(printf "> Link: (https://discord.com/channels/%d/%d/%d)" .Guild.ID .Channel.ID .Message.ID)
	"───────────────" }}

	{{ $_ := exec "remind" (str $time) $remText }}
	{{ $afterslice := split (exec "reminders") "\n"}}
	{{ range $afterslice}}{{if not (in $beforeslice .)}}{{with reFind `\*\*\d{8}\*\*` .}}{{$remID = slice . 2 10}}{{end}}{{end}}{{end}}
	{{ range $i,$_ := $keyWords}}{{if inFold $item .}}{{$emoji = index $pics $i}}{{end}}{{end}}

	
	{{/* sending Embed */}}
	{{ $embed := sdict
		"author" (sdict "name" (print "Auction #" (dbIncr 2 "auction_count" 1) " - Ongoing"))
		"thumbnail" (.User.AvatarURL "256" | sdict "url" )
		"description" $remText
		"color" 0x2f3136
		"timestamp" currentTime}}

	{{ $msgID := sendMessageRetID $reminder_channelID (complexMessage "content" (print .User.Mention ", You have started an auction. View reminders with `+reminders` command" ) "embed" $embed)}}
	{{ $newData := sdict "time" $time "item" $item "emoji" $emoji "content" $data.content "msgID" (str $msgID) "originalMsg" (str .Message.ID) "remID" (str $remID) "remText" $remText}}
	{{ dbSetExpire .User.ID "ongoing" $newData ($time.Seconds|toInt)}}
	{{ scheduleUniqueCC .CCID nil ($time.Seconds|toInt) .Message.ID $newData}}
	{{ addReactions "🏦" "🆎"}}
{{end}}


{{/* 🏦 */}}
{{if $bank_event}}
	{{sendMessage $reminder_channelID (joinStr "\n"
		(print (mentionRoleID $treasurer_role) "," .User.Mention "is requesting to transfer:")
		(print $data.emoji " **Item |** " $data.item)
		(print "⏱️ **Time |** %s"  $data.time)
			)
	}}
{{end}}

{{/* Closing Time */}}
{{if .ExecData}}
	{{$e := cembed "thumbnail" (.User.AvatarURL "256" | sdict "url" )
	"fields" (cslice 
		(sdict "name" "• Auction Details" "value" $data.content)
		(sdict "name" "• What should I do now?" "value" (printf "Change your auction status:\n```autohotkey\n- %-23s:  Open\n+ %-23s:  Closed```Announce Winner:\n```autohotkey\n- %-23s:\n+ %-23s:  None/@userMention```" "Status" "Status" "Winner" "Winner")))}}
	{{sendDM $e}}
	{{printf "%T,%v" $data.time $data.time}}
	{{/* Editing Closed Auction */}}
	{{ $remindEmbed := index ($fresh:= getMessage $reminder_channelID $data.msgID).Embeds 0|structToSdict}}
	{{ structToSdict $remindEmbed.Author | $remindEmbed.Set "Author"}}
	{{ reReplace "Ongoing" $remindEmbed.Author.Name "Closed" | $remindEmbed.Author.Set "Name"}}
	{{ $remindEmbed.Set "color" 0x53ddac}}
	{{ editMessage $fresh.ChannelID $fresh.ID (complexMessageEdit "embed" $remindEmbed) }}
	{{ deleteAllMessageReactions $auction_channelID $data.originalMsg}}
{{end}}
