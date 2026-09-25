# Prints the CHANGELOG.md section for version ver as GitHub release notes; joins wrapped bullet lines.
index($0, "## [" ver "]") == 1 { f = 1; print "## What's Changed"; next }
f && (/^## \[/ || /^\[[^]]+\]: /) { exit }
!f { next }
/^  [^ ]/ && buf != "" { sub(/^ +/, ""); buf = buf " " $0; next }
{ if (buf != "") print buf; buf = ""; if (/^- /) buf = $0; else print }
END { if (buf != "") print buf }
