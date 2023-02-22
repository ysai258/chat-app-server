# Use an official MySQL runtime as a parent image
FROM mysql:5.7

# Set the root password and create the chat database and user
ENV MYSQL_ROOT_PASSWORD=root
ENV MYSQL_DATABASE=chat
ENV MYSQL_USER=user
ENV MYSQL_PASSWORD=password

# Copy the schema.sql file to the /docker-entrypoint-initdb.d directory
COPY schema.sql /docker-entrypoint-initdb.d/

# Expose port 3306 for the MySQL server
EXPOSE 3306

# Install the Go runtime and set the working directory
RUN apt-get update && apt-get install -y golang
WORKDIR /app

# Copy the Go source code to the container
COPY . .

# Build the Go project
RUN cd cmd && go build -o app main.go

# Start the MySQL server and run the Go project
CMD ["sh", "-c", "/etc/init.d/mysql start && ./cmd/app"]

