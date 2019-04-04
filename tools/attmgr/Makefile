run: build
	./attmgr

trun: clean build
	./attmgr

build:
	go build --tags "sqlite_json"

sql:
	sqlite3 attestations.sqlite3

clean:
	-rm attestations.sqlite3

generate-entities:
	wv mke -o ns.ent --nopassphrase
	wv mke -o src.ent --nopassphrase
	wv mke -o dst1.ent --nopassphrase
	wv rtgrant --attester ns.ent --subject src.ent --indirections 3 -e 30d "wavemq:perm1,perm2,perm3@ns.ent/*"
	wv rtgrant --attester src.ent --subject dst1.ent --indirections 1 -e 30d "wavemq:perm1@ns.ent/abc/1" "wavemq:perm1@ns.ent/abc/2"
