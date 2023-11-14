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

.PHONY: clean
clean:
	rm -r $(SAVES)
