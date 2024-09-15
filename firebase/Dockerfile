# Use Node.js as the base image
FROM node:20

# Install Firebase CLI
RUN npm install -g firebase-tools
RUN apt-get update && apt-get install -y openjdk-17-jdk curl

# Set working directory
WORKDIR /firebase

# Copy Firebase configuration files
COPY firebase.json .
COPY .firebaserc .
COPY database.rules.json .

# Expose ports for Realtime Database and Auth emulators
EXPOSE 9000 9099

# Start Firebase emulators
CMD ["firebase", "emulators:start", "--only", "database,auth"]