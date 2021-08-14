# zfs-exporter

Use at your own peril. This exporter is far from stable. It works reasonably
well with an aggressive restart policy :) This crashes frequently and for
seemingly no reason at times. I'm aware that there are likely race bugs and
other strange interactions with the Go runtime that are causing these sporadic
panics and segmentation faults. Maybe I'll find and fix them, one day.

`go-libzfs` is based around github.com/bicomsystems/go-libzfs with a 
considerable portion of the libzfs wrapping logic rewritten. Credit goes to that
implementation for the framework and pointers on how to reference the largely
undocumented libzfs.
