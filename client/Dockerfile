# Use the Node.js 16 image as a base image
FROM node:18 as builder

# Copy the package.json and yarn.lock files
COPY package.json ./

# Install dependencies
RUN yarn install --frozen-lockfile

# Copy the rest of your project files into the image
COPY . .

# Build your Next.js application
RUN yarn build

# For the final image, use the Node.js 16 image as well
FROM node:18

# Set the working directory
WORKDIR /

# Copy the build from the previous stage
COPY --from=builder /.next ./.next
COPY --from=builder /node_modules ./node_modules
COPY --from=builder /public ./public
COPY --from=builder /package.json ./package.json

# Expose the port the app runs on
EXPOSE 3000

# Set the command to start the node server
CMD ["yarn", "start"]