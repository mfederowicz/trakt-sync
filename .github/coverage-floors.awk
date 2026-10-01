# Compares the per-package coverage in `go test -cover ./...` output with the floors file; exits 1 below a floor.
# Usage: awk -f .github/coverage-floors.awk .github/coverage-floors.txt test-output.txt
FNR == NR { if ($0 !~ /^#/ && NF == 2) floor[$1] = $2; next }
/coverage:/ {
	p = ""; c = ""
	for (i = 1; i <= NF; i++) { if ($i ~ /^github\.com\//) p = $i; if ($i == "coverage:") c = $(i+1) }
	sub("^github.com/mfederowicz/trakt-sync/?", "", p)
	sub("%", "", c)
	if (p != "") got[p] = c
}
END {
	n = 0
	for (p in floor) names[++n] = p
	# insertion sort keeps the report in a stable order without gawk extensions
	for (i = 2; i <= n; i++) { v = names[i]; for (j = i - 1; j >= 1 && names[j] > v; j--) names[j+1] = names[j]; names[j+1] = v }
	if (n == 0) { print "no coverage floors found"; exit 1 }
	bad = 0
	for (i = 1; i <= n; i++) {
		p = names[i]
		if (!(p in got)) { printf "FAIL %s: no coverage reported, floor is %s%%\n", p, floor[p]; bad = 1; continue }
		if (got[p] + 0 < floor[p] + 0) { printf "FAIL %s: %s%% is below the floor of %s%%\n", p, got[p], floor[p]; bad = 1; continue }
		printf "ok   %s: %s%% (floor %s%%)\n", p, got[p], floor[p]
	}
	exit bad
}
