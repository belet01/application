FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt /app/
RUN pip install --no-cache-dir -r requirements.txt
COPY . /app/
COPY templates/ /app/templates/
EXPOSE 5000
ENV FLASK_APP=user
CMD ["flask", "run", "--host=0.0.0.0", "--port=5000", "--debug"]