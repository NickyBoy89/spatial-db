SAVES := imperial-save
SAVES += skygrid-save
SAVES += witchcraft-save

all: $(SAVES)

imperial-save: compile
	./spatial-db load worldsave "saves/Imperialcity v14.1/region" --output "imperial-save"

skygrid-save: compile
	./spatial-db load worldsave "saves/SkyGrid/region" --output "skygrid-save"

witchcraft-save: compile
	./spatial-db load worldsave "saves/Witchcraft/region" --output "witchcraft-save"

.PHONY: compile
compile:
	CC=clang go build .

.PHONY: bench
bench: compile
	CC=clang GOEXPERIMENT=loopvar go test -bench .

.PHONY: clean
clean:
	rm -r $(SAVES)
