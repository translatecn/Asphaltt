TARGETSERVER = helloworld-server
TARGETCLIENT = helloworld-client
SRCDIRSERVER = server/*.c
SRCDIRCLIENT = client/*.c
BIN_DIR      = bin
INC_DIR      = -Iinclude
LIBRARY      = -luv
CC           = gcc

.PHONY: clean build

build:
	$(CC) $(INC_DIR) $(LIBRARY) -o $(BIN_DIR)/$(TARGETSERVER) $(SRCDIRSERVER)
	$(CC) $(INC_DIR) $(LIBRARY) -o $(BIN_DIR)/$(TARGETCLIENT) $(SRCDIRCLIENT)

clean:
	rm $(BIN_DIR)/*
