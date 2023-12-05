SAVES := imperial-save
SAVES += skygrid-save
SAVES += witchcraft-save

all: $(SAVES)

imperial-save:
	mkdir imperial-save
	./spatial-db load worldsave "saves/Imperialcity v14.1/region" --output "imperial-save"

skygrid-save:
	mkdir skygrid-save
	./spatial-db load worldsave "saves/SkyGrid/region" --output "skygrid-save"

witchcraft-save:
	mkdir witchcraft-save
	./spatial-db load worldsave "saves/Witchcraft/region" --output "witchcraft-save"

.PHONY: compile
compile:
	CC=clang go build .

.PHONY: bench
bench: compile
	CC=clang go test -bench . -benchtime=2s -count 10

.PHONY: clean
clean:
	rm -r $(SAVES)
