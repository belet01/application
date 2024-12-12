FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt /app/
RUN pip install -r requirements.txt
COPY . /app/
COPY templates/ /app/templates/
EXPOSE 5000
ENV FLASK_APP=user
ENTRYPOINT ["python", "main.py"]



