#!/bin/bash

TRGT=../mocks
SRC=../
VENDOR_PREF='github.com\/vwdilab\/flashlight-grasshopper\/vendor\/'

rm -r $TRGT
mkdir $TRGT

mockgen -package mock_grasshopper github.com/vwdilab/flashlight-grasshopper/grasshopper NewRelicFetcher,CloudFoundryFetcher | sed "s/$VENDOR_PREF//g"  > $TRGT/grasshopper_mocks.go
