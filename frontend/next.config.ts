import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  distDir: '.next',
  turbopack: {
    root: require('path').resolve(__dirname),
  },
  webpack: (config, { isServer }) => {
    config.resolve.alias['@'] = require('path').resolve(__dirname, 'src');
    return config;
  },
};

export default nextConfig;

// Only init OpenNext Cloudflare in development
if (process.env.NODE_ENV === 'development') {
  try {
    const { initOpenNextCloudflareForDev } = require('@opennextjs/cloudflare');
    initOpenNextCloudflareForDev();
  } catch (e) {
    // Ignore in production build
  }
}
