FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt /app/
RUN pip install -r requirements.txt

COPY . /app/
COPY templates/ /app/templates/
COPY static/ /app/static/
EXPOSE 5000
CMD ["flask", "run", "--host=0.0.0.0", "--port=5000"]
