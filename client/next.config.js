/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  env: {
    baseURL: process.env.NEXT_PUBLIC_BASE_URL,
  },
};

module.exports = nextConfig;
