for i in {1..10000}; do 
	for i in {1..100}; do 
		curl -i localhost:8000/work --data key=K$i -G; 
	done
done