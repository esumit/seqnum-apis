#!/bin/bash

RANDOM_GETPOST=600

get_seqnum(){
  echo
  echo "GET <<<"
  echo "URL: http://localhost:9000/seqnum"
  echo
  result=$(curl -s -i --location --request GET --header 'Content-Type: application/json' 'http://localhost:9000/seqnum')
  echo
  echo "Response:"
  echo
  echo $result
  echo "-----------------------------------------------"
}

for ((i = 1; i < RANDOM_GETPOST; ++i)); do
	echo
	echo "Genreate Seqnum GET Request : $i "
	echo "Press [CTRL+C] to stop.."
  get_seqnum
  sleep 3
done
