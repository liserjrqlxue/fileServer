<html>
<head>
	<title>MP4 Play</title>
</head>
<body>

<input type="hidden" name="token" value="{{.Token}}"/>
<!-- first try HTML5 playback: if serving as XML, expand `controls` to `controls="controls"` and autoplay likewise -->
<!-- warning: playback does not work on iOS3 if you include the poster attribute! fixed in iOS4.0 -->
<video width="640" height="360" controls>
	<!-- MP4 must be first for iPad! -->
	<source src="{{.Src}}" type="video/mp4" /><!-- Safari / iOS video    -->
	<source src="{{.Src}}" type="video/ogg" /><!-- Firefox / Opera / Chrome10 -->
	<!-- fallback to Flash: -->
	<object width="640" height="360" type="application/x-shockwave-flash" data="{{.Src}}">
		<!-- Firefox uses the `data` attribute above, IE/Safari uses the param below -->
		<param name="movie" value="{{.Src}}" />
		<param name="flashvars" value="controlbar=over&amp;image=__POSTER__.JPG&amp;file={{.Src}}" />
		<!-- fallback image. note the title field below, put the title of the video there -->
		<img src="__VIDEO__.JPG" width="640" height="360" alt="__TITLE__"
		     title="No video playback capabilities, please download the video below" />
	</object>
</video>
<!-- you *must* offer a download link as they may be able to play the file locally. customise this bit all you want -->
<p>	<strong>Download Video:</strong>
	Closed Format:	<a href="{{.Src}}.MP4">"MP4"</a>
	Open Format:	<a href="__VIDEO__.OGV">"Ogg"</a>
</p>

</body>
