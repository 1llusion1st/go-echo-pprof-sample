collect_pprofile:
	curl http://localhost:8080/debug/pprof/profile?debug=1 > profile.prof

collect_allocs:
	curl http://localhost:8080/debug/pprof/allocs?debug=1 > allocs.prof

collect_heap:
	curl http://localhost:8080/debug/pprof/heap?debug=1 > heap.prof

_get_exe:
	curl http://localhost:8080/debug/pprof/cmdline > executable.path

show_pprofile: _get_exe
	go tool pprof -web 8081 $(shell cat executable.path) profile.prof

show_allocs: _get_exe
	go tool pprof -web 8081 $(shell cat executable.path) allocs.prof

show_heap: _get_exe
	go tool pprof -web 8081 $(shell cat executable.path) heap.prof
