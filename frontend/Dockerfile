# Step 1: Build the Next.js app
FROM node:18-alpine AS build

# Set working directory
WORKDIR /app

# Install dependencies
COPY frontend/package*.json ./
RUN npm install

# Copy the whole Next.js app
COPY frontend/ ./

# Build the Next.js app (creates optimized production build)
RUN npm run build

# Step 2: Serve the app with Next.js
FROM node:18-alpine AS production

WORKDIR /app

# Copy the build from the build stage
COPY --from=build /app /app

# Install only production dependencies
RUN npm install --production

# Expose the default Next.js port
EXPOSE 3000

# Start the Next.js server
CMD ["npm", "start"]
