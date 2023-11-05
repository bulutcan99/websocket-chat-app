/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  env: {
    baseURL: process.env.BASE_URL,
  },
};

module.exports = nextConfig;
