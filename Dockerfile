FROM python:3.6.3
WORKDIR /python/src/github.com/callstats-io/ai-decision
ADD requirements.txt .
RUN pip install -r requirements.txt
ADD main.py .
ADD backfill.pickle .
ADD src/ src/
ADD protos/ protos/
ADD service/gen/protos/ service/gen/protos/
ADD scripts/ scripts/
CMD python main.py --manual_date $MANUAL_DATE --unsuppress $UNSUPPRESS
