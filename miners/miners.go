{{/* 
      Recommended Trigger: \A(?:\+|<@!?204255221017214977>)\s*(?:miners?)(?: +|\z)
      Recommended Trigger Type: Regex

      Usage: +miners <server:Text> <gate:Text> <mob:Whole Number>

    	Copyright (c): zen | ゼン#0008; 2021
	    License: MIT
      Repository: https://github.com/z3nn13/System-Dungeons
*/}}

{{/* Configuration Starts */}}
{{ $mineRole := 759698366589435914}}
{{/*{{ $mineRole := 840493186535129088 }} */}}
{{ $servers := cslice "main" "c-main" "support" "c-support"}}
{{ $mobs := cslice 10 5 1}}
{{ $guilds := sdict
  "guild" (sdict

  "main" 592385635143122946
  "c-main" 592385635143122946
  "support" 617697836674449434
  "c-support" 617697836674449434)

  "channels" (sdict 
  "main" 662900098283470858
  "mainForced" 613931390282629131
  "c-main" 965679980463927347
  "c-mainForced" 964792877366534144

  "support" 706441024867926116
  "supportForced" 969184547164008498
  "c-support" 964777431682646097
  "c-supportForced" 654366614355050516 )}}

{{ $thumbnail := "https://media.discordapp.net/attachments/662115940296687646/662389317041258506/Mana_Crystal.png"}}
{{ $avatar := sdict 
  "main" "https://cdn.discordapp.com/icons/592385635143122946/15de84e3c311687b251b520e99fd2703.webp"
  "c-main" "https://cdn.discordapp.com/icons/592385635143122946/15de84e3c311687b251b520e99fd2703.webp"
  "support" "https://cdn.discordapp.com/attachments/710699457867546694/959830273749569586/2020.10.13_22.14.04.gif"
  "c-support" "https://cdn.discordapp.com/icons/617697836674449434/a_09ec0a0b848a328aac55ba7295576d26.gif"
}}

{{/* Don't recommend changing this*/}}
{{ $gates := cslice "absolute" "rulers" "antares" "monarch" "s" "sred" "a" "ared" "bred" "b" "cred" "c"}}
{{ $aliases := cslice "abs" "ruler" "ant" "m"}}
{{ $gateEmojis := sdict
  "blue" "<a:Blue_Gate:824690207676301372>"
  "red" "<a:Red_Gate:824691553682456576>"
  "black" "<:Black_Gate:942064710424088596>"
  "white" "<:White_Gate:961479511826919504>"
}}
{{ $colors := sdict
"main" 0xcea1f7 
"c-main" 0xcea1f7
"support" 0xffbddb
"c-support" 0xffbddb}}
{{$mineCooldown := 15}}{{/* seconds */}}
{{/* Configuration Ends */}}

{{/* Global Variables */}}
{{ $help := sdict "title" "Miner/miners" "description" (joinStr "\n" "```md"
  "Miners <server:Text> <gate:Text> <mob:Whole Number>``````"
  "[-forced forced:Switch - breaks in dungeon channels]```") }}
{{$server := ""}}{{$gate := ""}}{{$mob := ""}}
{{$newArgs := ""}}{{$forced := false}}
{{$errMsg := ""}}{{$err := false}}{{$allowed := ""}}


{{/* Forced switch */}}
{{if reFind `\s*(?:-)?(?i:forced)(?:\s+|\z)` .StrippedMsg}}
    {{$newMsg := reReplace `\s*(?:-)?(?i:forced)` .StrippedMsg ""}} 
    {{with $newMsg}}{{$newArgs = split . " "}}{{end}}
    {{$forced = true}}
{{end}}

{{/* Parse input */}}
{{$CmdArgs := or $newArgs .CmdArgs}}
{{with $CmdArgs}}
    {{if ge (len .) 3}}
            {{$server = index . 0 | lower }}
            {{$gate = index . 1 | lower }}
            {{$mob = index . 2 | lower | toInt }}
    {{else}}
        {{$errMsg = "Error : Insufficent Arguments Provided"}}{{$err = true}}
    {{end}}
{{else}}
    {{$err = true}}{{sendMessage nil (cembed $help)}}
{{end}}

{{if not $err}}
    {{/*   */}}
    {{if not (in $servers $server)}}{{$errMsg = printf "Unknown Server : %q" $server}}{{$allowed = json $servers}}
    {{else if not (or (in $gates $gate) (in $aliases $gate))}}{{$errMsg = printf "Unknown Gate Rank: %q" $gate}}{{$allowed = json $gates}}
    {{else if not (in $mobs $mob)}}{{$errMsg = printf "Invalid Mob Count: %q" (str $mob)}}{{$allowed = json $mobs}}
    {{end}}
{{end}}

{{if and (not $err) (not $errMsg)}}
  {{/* Cooldown Part */}}
  {{ with dbGet 2 "mineCD" }}
    {{sendMessage nil (cembed "color" 0xffe4d9 "description" (print "🔒 Command Locked: Reopening in `" (.ExpiresAt.Sub currentTime | humanizeDurationSeconds) "`"))}}
  {{ else }}

    {{/* Alias part */}}
    {{if in $aliases $gate}}{{range $i,$_ := $aliases}}{{if eq . $gate}}{{$gate = index $gates $i}}{{end}}{{end}}{{end}}
    
    {{/* Flags */}}
    {{$is_Sunday := eq (str currentTime.Weekday) "Sunday"}}
    {{/* {{$is_Sunday = true}} */}}
    {{$is_Red := or $is_Sunday (in $gate "red") false}}
    {{$is_Custom := in $server "c-"}}
    {{if $is_Custom}}{{$server = slice $server 2}}{{end}}

    {{/* Default */}}
    {{$gateEmoji := $gateEmojis.blue}}
    {{ if $is_Red}}
      {{if in $gate "red"}}{{$gate = len $gate|add -3|slice $gate 0}}{{end}}
      {{ $gateEmoji = $gateEmojis.red }}
      {{ $is_Red = " `[Red]`" }}
    {{ else if eq $gate "antares" "monarch"}}
        {{ $gateEmoji = $gateEmojis.black }}    
    {{ else if eq $gate "absolute" "rulers"}}
        {{ if eq $gate "absolute"}}{{$gate = "absolute being"}}{{end }}
        {{ $gateEmoji = $gateEmojis.white }}
    {{ end }}

    {{/* Embed Generation */}}
    {{$note := or (and $is_Sunday "<:pink_dot:962206254527303730> **Note:** All gates are red on Sunday") "<:pink_dot:962206254527303730> **[** Tip: Be sure to read 📝**[mine rules](https://discord.com/channels/617697836674449434/706441024867926116/832898336297975809) ]**"}}
    {{$channel_link := printf "https://discord.com/channels/%d/%d" ($guilds.guild.Get $server) (or (and $forced (print $server "Forced")) $server|$guilds.channels.Get)}}
    {{ $embed := cembed 
      "thumbnail" (sdict "url" $thumbnail)
      "color" ($colors.Get $server)
      "author" (sdict "name" (print (title $server) " Server") "icon_url" ($avatar.Get $server))
      
      "title" (joinStr "\n" 
      (printf "> • <a:abell:823825677747748874> **%d %s left**" $mob (or (and (eq $mob 1) "mob") "mobs"))
      (printf "> • **%s %s break%s**" $gateEmoji (title (print (and $is_Custom "Custom "|str) $gate)) ($is_Red|str)))
      
      "description" (joinStr "\n"
      (printf "\u200b")
      $note
      (printf "<a:moon_sparkle:1030310798855262240> **[[Click To Jump To The Moon](%s)]**" $channel_link))
      }}
    {{ sendMessageNoEscape nil (complexMessage "content" (mentionRoleID $mineRole) "embed" $embed)}}
    {{ dbSetExpire 2 "mineCD" true $mineCooldown}}
    {{  $rep := execAdmin "+" .User}}
    {{ sendMessage nil (joinStr "\n" 
      (print "Thanks for pinging <:birblove:858995265927249940> " (or (and (in $rep "Error:") "(Rep on cooldown)") $rep)))
      }}

  {{end}}
{{end}}

{{/* Error output */}}
{{with $errMsg}}
  {{with $allowed}}
    {{$help.Set "fields" (cslice (sdict "name" "• Acceptable Arguments" "value" (print "```elm\n" . "```")))}}
  {{end}}
  {{sendMessage nil (complexMessage "content" . "embed" $help)}}
{{end}}