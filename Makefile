.PHONY: doky func stress stop

doky:
	docker build -t testy .
	docker run -p 5000:5000 --rm --name testy -t testy

stop:
	docker stop testy

func:
	./technopark-dbms-forum func -u http://localhost:5000/api -r report.html

stress:
	./technopark-dbms-forum fill --url=http://localhost:5000/api --timeout=900
	./technopark-dbms-forum perf --url=http://localhost:5000/api --duration=600 --step=60
