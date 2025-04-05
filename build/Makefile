
all:
	cd db && make;
	cd kafka && make;

compose-up:
	cd db && make compose-up;

migrations-up:
	cd db && make goose-up;

compose-stop:
	cd db && make compose-down;

migrations-down:
	cd db && make goose-down-all;