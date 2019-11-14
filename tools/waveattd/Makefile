run: build
	./waveattd

trun: build
	./waveattd

build:
	go build --tags "sqlite_json"

install:
	go install --tags "sqlite_json"

sql:
	sqlite3 attestations.sqlite3

#clean:
#	-rm -f attestations.sqlite3

generate-entities: build
	# clean up from last time
	rm -f *.pem src.db
	wv mke -o ns.ent --nopassphrase
	wv mke -o src.ent --nopassphrase
	wv mke -o dst1.ent --nopassphrase
	wv name --attester src.ent dst1.ent dst1
	wv name --attester src.ent ns.ent testns
	wv rtgrant --attester ns.ent --subject src.ent --indirections 3 -e 30d "wavemq:perm1,perm2,perm3@ns.ent/*"
	wv rtgrant --attester src.ent --subject dst1 --indirections 1 -e 30d "wavemq:perm1@testns/abc/1" "wavemq:perm1@testns/abc/2"
	WAVE_DEFAULT_ENTITY=src.ent WAVEATTD_DB_LOCATION=src.db ./waveattd
