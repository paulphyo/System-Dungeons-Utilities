{{/*
	Recommended Trigger: Reaction Added
	Recommended Trigger Type: Reaction

	Usage: Execute auction start/stop events based on reactions

	Copyright (c): zen | ゼン#0008; 2021
	License: MIT
	Repository: https://github.com/z3nn13/System-Dungeons
*/}}


{{/* Configuration Variables Start */}}
{{$reminder_channelID := "819251873746780241"}}
{{/* Configuration Variables End */}}

{{/* Actual Code */}}
{{$emoji_name := .Reaction.Emoji.Name}}
{{$ex := or (and (reFind "a_" .Guild.Icon) "gif" ) "png" }}
{{$icon := print "https://cdn.discordapp.com/icons/" .Guild.ID "/" .Guild.Icon "." $ex "?size=1024" }}


{{$z := sdict}}
{{$slice := cslice}}
{{$title := cslice}}
{{$titleL := cslice}}

{{$t := ""}}
{{$duration := ""}}
{{$convertedDuration := ""}}
{{$timeLeft := ""}}
{{$formattedTimeLeft := ""}}
{{$mid := ""}}
{{$embed1 := sdict
"author" (sdict
"name" (print .ReactionMessage.Author.Username " has started an auction!")
"icon_url" (.ReactionMessage.Author.AvatarURL "1024"))
"color" 0x53ddac
"timestamp" currentTime}}

{{if .ExecData }}
    {{$time := toDuration .ExecData.duration}}
    {{$timeLeft2 := currentTime.Add $time}}
    {{$item := .ExecData.itemName}}
    {{$channel1 := 750295348014219305}}
    {{$channel2 := 819251873746780241}}
    {{$mid := .ExecData.message}}
    {{$rMsgID := .ExecData.rmsgID}}
    {{$refresh := getMessage $channel1 $rMsgID}}
    {{$embed2 := sdict
    "author" (sdict
    "name" (print .User.Username "'s auction has ended")
    "icon_url" (.User.AvatarURL "512"))
    "timestamp" currentTime
    "color" 0xff5a54}}
    {{$embed2.Set "fields" (cslice 
    (sdict "name" "Auction" "value" (print "[Click Here](https://discord.com/channels/" .Guild.ID "/" .Channel.ID "/" $refresh.ID ")") "inline" true)
    (sdict "name" "Channel" "value" (print "<#" $channel1 ">") "inline" true)
    (sdict "name" "Duration" "value" (print "<a:hg:832130433033699338> " (humanizeDurationSeconds $time)) "inline" true)
    (sdict "name" "End Date" "value" ($timeLeft2.Format "🗓️ 02 Jan 2006\n🕐 3:04 PM MST") "inline" true)
    (sdict "name" "Item" "value" $item "inline" true))}}
    {{editMessage $channel2 $mid (complexMessageEdit "embed" (cembed $embed2))}}
    {{$embed2.Set "fields" (cslice 
    (sdict "name" "Auction Info" "value" $refresh.Content "inline" true)
    (sdict "name" "Auction Link" "value" (print "[Jump to](https://discord.com/channels/" .Guild.ID "/" .Channel.ID "/" $refresh.ID ")") "inline" false)
    (sdict "name" "Auction Link" "value" (print "<#" .Channel.ID ">") "inline" true)
    (sdict "name" "Notice" "value" (print "<a:a_Check_Mark:832292968319942666> Update your auction status to closed in the next **24 hours**\n<a:a_Check_Mark:832292968319942666> Failure to do this = __warn/mute__\n<a:a_Check_Mark:832292968319942666> Ask in <#762942886991757362> if you have any questions") "inline" false))}}
    {{sendDM (cembed $embed2)}}
    {{deleteAllMessageReactions nil $refresh.ID}}
{{else}}
{{if and (eq $emoji_name "⏰" "🗑️" "🏦") (not (eq .ReactionMessage.Author.ID .User.ID))}}
    {{deleteMessageReaction nil .ReactionMessage.ID .User.ID $emoji_name}}
{{else}}
{{if (eq $emoji_name "🗑️")}}
    {{deleteMessage nil .ReactionMessage.ID 0}}
{{else if (eq $emoji_name "🏦")}}
    {{deleteMessageReaction nil .ReactionMessage.ID .User.ID $emoji_name}}
    {{$slice = reFindAllSubmatches `(?m)\A(?:\x60{3}(?:\w+)?)?((?:\n)?Item[\w\s+.-]+)\:([\w\s+!@#$%^\x27&*\/\x5b\x5c\x28\x29 _';,.-]+)$\n^([\w\s+.-]+)\:([\w\s.-]+)$\n(?:^([\w\s+.-]+)\:([\w\s]+)$\n)?^([\w\s.]+)\:([\w\s.]+)$\n^([\w\s.-]+)\:([\w\s+.-]+)$\n^([\w\s+.-]+)\:([\w\s+\x60.-]+)$(?:\n\x60{3})?(?:\n^([\w\s+.-]+)\:)?(?:\n\x60{3})?(?:\s+|\z)` .ReactionMessage.Content}}
    {{$embed1.Set "author" (sdict "name" (print .User.Username "(" .User.ID ")") "icon_url" (.User.AvatarURL "512"))}}
    {{$embed1.Set "color" 0xdfa1ff}}
    {{$embed1.Set "fields" (cslice (sdict "name" "Item" "value" (index $slice 0 2) "inline" true) (sdict "name" "Duration" "value" (print "<a:hg:832130433033699338>" (index $slice 0 10)) "inline" true) (sdict "name" "Link" "value" (print "[Click Here](https://discord.com/channels/" .Guild.ID "/" .Channel.ID "/" .ReactionMessage.ID ")") "inline" true))}}
    {{$embed1.Set "thumbnail" (sdict "url" "https://media.discordapp.net/attachments/806003815354335263/832838894181023755/image0.png?width=498&height=498")}}
    {{$embed1.Del "footer"}}
    {{$embed1.Del "timestamp"}}
    {{sendMessageNoEscape $reminder_channelID (complexMessage "content" (print "<@&818732994159181835>, " .User.Mention " has requested a transaction!") "embed" (cembed $embed1))}}
    {{deleteAllMessageReactions nil .ReactionMessage.ID}}
{{else if eq $emoji_name "⏰"}}
    {{deleteAllMessageReactions nil .ReactionMessage.ID}}
    {{addReactions "🏦"}}
    {{$slice = reFindAllSubmatches `(?m)\A(?:\x60{3}(?:\w+)?)?((?:\n)?Item[\w\s+.-]+)\:([\w\s+!@#$%^\x27&*\/\x5b\x5c\x28\x29 _';,.-]+)$\n^([\w\s+.-]+)\:([\w\s.-]+)$\n(?:^([\w\s+.-]+)\:([\w\s]+)$\n)?^([\w\s.]+)\:([\w\s.]+)$\n^([\w\s.-]+)\:([\w\s+.-]+)$\n^([\w\s+.-]+)\:([\w\s+\x60.-]+)$(?:\n\x60{3})?(?:\n^([\w\s+.-]+)\:)?(?:\n\x60{3})?(?:\s+|\z)` .ReactionMessage.Content}}
    {{$duration = (index $slice 0 10)}}
    {{$convertedDuration = toDuration $duration}}
    {{$t = currentTime.Add $convertedDuration}}
    {{$timeLeft = $t.Sub currentTime}}
    {{$formattedTimeLeft = humanizeDurationSeconds $timeLeft}}
    {{$embed1.Set "fields" (cslice 
    (sdict "name" "Auction" "value" (print "[Click Here](https://discord.com/channels/" .Guild.ID "/" .Channel.ID "/" .ReactionMessage.ID ")") "inline" true)
    (sdict "name" "Channel" "value" (print "<#" .Channel.ID ">") "inline" true)
    (sdict "name" "Duration" "value" (print "<a:hg:832130433033699338> " (humanizeDurationSeconds $convertedDuration)) "inline" true)
    (sdict "name" "End Date" "value" ($t.Format "🗓️ 02 Jan 2006\n🕐 3:04 PM MST") "inline" true)
    (sdict "name" "Item" "value" (index $slice 0 2) "inline" true))}}
    {{$mid = sendMessageRetID $reminder_channelID (complexMessage
    "content" (print .User.Mention " (You can view reminders with the reminders command)") 
    "embed" (cembed ($embed1)))}}
{{$z.Set "duration" $convertedDuration}}
{{$z.Set "message" $mid}}
{{$z.Set "rmsgID" .ReactionMessage.ID}}
{{$z.Set "itemName" (index $slice 0 2)}}
{{$reminderText := print "\n ───※ ·❆· ※─── \nYour auction has ended! Check your DMs for detailed info.\n Auction Link: (https://discord.com/channels/" .Guild.ID "/" .Channel.ID "/" .ReactionMessage.ID ")\n **───※ ·❆· ※───**"}}
{{$silent := exec "reminder" (print $convertedDuration) $reminderText}}
{{execCC .CCID nil $convertedDuration.Seconds $z}}
{{else}}
{{end}}    
{{end}}
{{end}}

