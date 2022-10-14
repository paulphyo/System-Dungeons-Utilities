{{/*
	Recommended Trigger: (?:\x60{3}\S*)?(?:[^:]+:[^\n]+\n){3,}
	Recommended Trigger Type: Regex

	Usage: Listens a format in #dungeon-auction and raise reactions

	Copyright (c): zen | ゼン#0008; 2022
	License: MIT
	Repository: https://github.com/z3nn13/System-Dungeons
*/}}

{{$mod_Role = 699868880188866631}}

{{define "error"}}
		{{$embed1 := cembed
			"color" 0x8a2be2
			"author" (sdict "name" "Your aucton has been automatically deleted!")
			"icon_url" (.User.AvatarURL "1024")
			"thumbnail" (sdict "url" "https://media.discordapp.net/attachments/806003815354335263/832838894181023755/image0.png")
			"fields" (cslice
					(sdict "name" "• What you tried" "value" .TemplateArgs.content)
					(sdict "name" "• Reasons" "value" (joinStr "\n" .TemplateArgs.errors.StringSlice))
					(sdict "name" "• Help" "value"
					(joinStr "\n" 
					"> To Check the **Format**, [Click Here](https://discord.com/channels/617697836674449434/750295348014219305/750305781781364737)"
					"> To Read the **Auction Rules**, [Click Here](https://discord.com/channels/617697836674449434/750295348014219305/819109901473284096)"
					"> Ask For Help in <#762942886991757362>") "inline" true))}}
		{{sendDM $embed1}}
{{end}}

{{/*	Sample
```autohotkey
Item Name + Quantity    : 10 x Blessed Random Box
Starting Price          : 189k Mana Crystal
Auto-buy                : 199k Mana Crystal
Min. Increase           : 5k Mana Crystal
Time Period             : 12h
Status                  : Open
```
**Winner**                           :
*/}}


{{/* Global Variables */}}
{{$input := reReplace `\x60[\S*_]*|\s{2,}|[*_]+` .Message.Content ""}}{{/* Trimming stylish stuff */}}
{{$lines := split $input "\n"}}
{{$visuals := sdict}}{{$auction := sdict}}{{$keys := cslice}}
{{$order := split "abcdefghijklmnopqrstuvwxyz" ""}}
{{/* TODO: 
		Don't trim value, only keys
	*/}}

{{/* Fetching keys from input */}}
{{range $i,$_ := $lines}}
	{{if .}}
	{{- $split := split . ":"}}
	{{- $key := index $split 0 | lower}}{{$val := index $split 1}}
	{{- $visuals.Set (print (index $order $i) ($key := index $split 0 | lower)) ($val := index $split 1)}}
	{{- $auction.Set $key $val}}
	{{- $keys = $keys.Append $key}}{{end}}
{{- end}}

{{/* Checking input format */}}
{{$err := false}}{{$wrong_key := ""}}
{{$form_check := cslice "item name + quantity" "starting price" "min. increase" "status" "time period" "winner"}}
{{$spell_check:= cslice "auto-buy" "autobuy" | $form_check.AppendSlice}}
{{$errorMsgs  := cslice}}
{{range $keys}}{{if not (in $spell_check .)}}{{$errorMsgs = $errorMsgs.Append (printf "<a:no:832258616127389736> Unknown Key: %q" (title .))}}{{end}}{{end}}
{{range $form_check}}{{ if not (in $keys .)}}{{$errorMsgs = $errorMsgs.Append (printf "<a:no:832258616127389736> %q missing" (title .))}}{{end}}{{end}}


{{/* Time validity */}}
{{ if $errorMsgs }}
		{{ with reFindAllSubmatches `(\d+) ?(\w+)` $auction.Get "time period" }}
				{{ $duration := index . 0 0 }}
				{{ with toDuration $duration }}
						{{ $auction.Set "time period" $duration }}
				{{ else }}
						{{ $errorMsgs = $errorMsgs.Append (printf "<a:no:832258616127389736> Invalid Time Argument: %q" $duration) }}
				{{ end }}
		{{ else }}
			{{ $errorMsgs = $errorMsgs.Append (printf "<a:no:832258616127389736> Invalid Time Argument: %q\nExamples: `12h`,`1 hour`,`2 days`" $auction.Get "time period")}}
		{{ end }}
{{ end }}

{{/* Prettify output  */}}
{{$content := "```autohotkey\n"}}
{{range $k,$v := $visuals}}
	{{- $content = printf "%s%-25s: %s\n" $content (slice $k 1|title) (.|title) -}}
{{end}}
{{$content = print $content "```"}}


{{ if $errorMsgs }}
		{{ deleteTrigger 1 }}
		{{ if hasRoleID $mod_Role }}
				{{ sendDM "wow a mod breaking an auction rule, bonk" }}
		{{ else }}
				{{ execAdmin "warn" .User.ID "Invalid Auction Format" }}
		{{ end }}
		{{ sendTemplateDM "error" (sdict "errors" $errorMsgs "content" $content) }}
{{ else }}
		{{ addReactions "⏰" "🗑️" }}
		{{ dbSetExpire .Message.ID "can_auction" (sdict "auction" $auction "content" $content) 30 }}
{{ end }}