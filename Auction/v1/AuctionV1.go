{{/*
	Recommended Trigger: (?m)\A(?:\x60{3}(?:\w+)?)?((?:\n)?Item[\w\s+.-]+)\:([\w\s+!@#$%^\x27&*\/\x5b\x5c\x28\x29 _';,.-]+)$\n^([\w\s+.-]+)\:([\w\s.-]+)$\n(?:^([\w\s+.-]+)\:([\w\s]+)$\n)?^([\w\s.]+)\:([\w\s.]+)$\n^([\w\s.-]+)\:([\w\s+.-]+)$\n^([\w\s+.-]+)\:([\w\s+\x60.-]+)$(?:\n\x60{3})?(?:\n^([\w\s+.-]+)\:)?(?:\n\x60{3})?(?:\s+|\z)
	Recommended Trigger Type: Regex

	Usage: Listens a format in #dungeon-auction and raise reactions

	Copyright (c): zen | ゼン#0008; 2021
	License: MIT
	Repository: https://github.com/z3nn13/System-Dungeons
*/}}


{{/* Code */}}
{{$slice := reFindAllSubmatches `(?m)\A(?:\x60{3}(?:\w+)?)?((?:\n)?Item[\w\s+.-]+)\:([\w\s+!@#$%^\x27&*\/\x5b\x5c\x28\x29 _';,.-]+)$\n^([\w\s+.-]+)\:([\w\s.-]+)$\n(?:^([\w\s+.-]+)\:([\w\s]+)$\n)?^([\w\s.]+)\:([\w\s.]+)$\n^([\w\s.-]+)\:([\w\s+.-]+)$\n^([\w\s+.-]+)\:([\w\s+\x60.-]+)$(?:\n\x60{3})?(?:\n^([\w\s+.-]+)\:)?(?:\n\x60{3})?(?:\s+|\z)` .Message.Content}}
{{$embed1 := sdict
	"author" (sdict
	"name" (print "Warning: " .User.String " has been warned!")
	"icon_url" (.User.AvatarURL "1024"))
	"fields" (cslice
			(sdict "name" "Your auction has also been __automatically deleted__! <a:abell:823825677747748874>" "value" "• This is due to one of the following:" "inline" false)
			(sdict "name" "Reasons" "value" "<a:no:832258616127389736> Wrong format/spelling\n<a:no:832258616127389736> Against auction rules\n<a:no:832258616127389736> Click these links to recheck rules & format." "inline" false)
			(sdict "name" "Format" "value" "[Click Here](https://discord.com/channels/617697836674449434/750295348014219305/750305781781364737)" "inline" true)
			(sdict "name" "Rules" "value" "[Click Here](https://discord.com/channels/617697836674449434/750295348014219305/819109901473284096)" "inline" true)
			(sdict "name" "Help" "value" "<a:a_Check_Mark:832292968319942666> Request help in <#762942886991757362>\n<a:a_Check_Mark:832292968319942666> DM or ping an online staff member if you're experiencing repetitive problems!" "inline" false))
	"image" (sdict "url" "https://media.discordapp.net/attachments/806003815354335263/832838894181023755/image0.png?width=498&height=498")
	"color" 0x8a2be2}}
	
{{/* Removing whitespace */}}
{{$title := cslice}}
{{$matched := true}}
{{if not $slice}}
	{{$matched = false}}
{{else}}
{{range seq 1 (len (index $slice 0))}}
	{{- if (eq (toInt (mod . 2)) 0) -}}
	{{- else -}}
	{{- $title = $title.AppendSlice (reFindAll `(\S+)` (index $slice 0 .)) -}}
	{{end -}}
{{end}}
{{end}}

{{/*Converting to lower case */}}
{{$titleL := cslice}}
{{range seq 0 (len ($title))}}
	{{- $titleL = $titleL.Append ((index $title .)|lower) -}}
{{end}}

{{/*The Matching Slice*/}}
{{$fSlice := cslice "name" "+" "quantity" "starting" "price" "min." "increase" "time" "status" "winner"}}
{{$optional := cslice "item" "autobuy" "auto-buy" "period"}}
{{$sSlice := $fSlice.AppendSlice $optional}}
{{$format := true}}
{{$spelling := true}}

{{/* Checking spelling and format*/}}
{{range $fSlice}}
	{{- if not (in $titleL .)}}
	{{- $format = false -}}
	{{- else -}}
	{{end -}}
{{end}}
{{range $titleL}}{{- if not (in $sSlice .)}}{{- $spelling = false -}}{{- else -}}{{end -}}{{end}}
{{if or (eq $format false) (eq $spelling false) (eq $matched false)}}
	{{deleteMessage nil .Message.ID 0}}
	{{execAdmin "warn" .User.ID "Wrong auction format/spelling"}}
	{{sendDM (cembed $embed1)}}
{{else}}
	{{addReactions "⏰" "🗑️"}} {{sleep 15}}
	{{deleteAllMessageReactions nil .Message.ID}}
{{end}}
