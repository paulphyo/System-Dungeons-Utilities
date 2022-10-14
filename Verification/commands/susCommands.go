{{/*
        Trigger: \A(?:\+|<@!?204255221017214977>\s*)(?:sus|ns|notsus)(?: +|\z)
        Trigger Type: Regex

        Usage:
        +sus <User:Mention/ID> to add Suspicious Role + remove Human Role
        +notsus <User:Mention/ID> to add Human Role + remove Suspicious Role

		Copyright (c): zen | ゼン#0008; 2021
	    License: MIT
	    Repository: https://github.com/z3nn13/System-Dungeons
*/}}


{{/* Customize here if role IDs are changed */}}
{{ $prefix := "+"}}
{{ $roles := sdict 
    "sus" 764049927445020674
    "human" 645226278902956033 }}


{{/* Actual code. Don't touch */}}
{{ $cond := eq (lower .Cmd) (print $prefix "sus")}}
{{ $susID := or (and $cond $roles.sus) $roles.human }}
{{ $humanID := or (and $cond $roles.human) $roles.sus}}
{{ $err := ""}}

{{/* Parsing input */}}
{{ $args := parseArgs 1 (print "\x60\x60\x60" .Cmd " <User:Mention/ID>\x60\x60\x60")
    (carg "userid" "Member to take action on")
    (carg "string" "Message to ignore") }}
{{ $target := $args.Get 0|getMember }}


{{/* Conditions input */}}
{{ if not $target }}
    {{ $err = "Member not found" }}
{{ else }}
    {{ $target = $target.User}}
    {{ if eq $target.ID .User.ID }}
        {{ $err = "You cannot use this command upon yourself"}}
    {{ else }}
        {{ giveRoleID $target.ID $susID}}
        {{ takeRoleID $target.ID $humanID}}
        {{ addReactions ":Check:723793166196801567" }}
        {{ $embed := cembed
                "author" (sdict "name" (print .User.String "(ID " .User.ID ")") "icon_url" (.User.AvatarURL "256"))
                "color" 0x49ED47
                "description" (printf "<:check_mark:838811489581137961> **Finished Removing Role:** <@&%d>\n<:check_mark:838811489581137961> **Finished Adding Role: **<@&%d>\n🌙 **Member:** %s (ID %d)" $humanID $susID $target.String $target.ID)
                "thumbnail" (sdict "url" ($target.AvatarURL "256"))
        }}
        {{ sendMessage nil $embed}}
    {{ end }}
{{ end }}

{{ with $err }}
    {{- addReactions "❌" }}
    {{- with sendMessageRetID nil . }}{{ deleteMessage nil . 5 }}{{ end }}
{{ end }}