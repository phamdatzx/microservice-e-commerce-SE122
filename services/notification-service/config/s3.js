/**
 * AWS S3 Configuration
 * Initialize S3 client for file uploads
 */

import { S3Client } from '@aws-sdk/client-s3';
import dotenv from 'dotenv';

dotenv.config();

let s3Client = null;
let s3BucketName = null;

/**
 * Initialize S3 client with credentials from environment
 */
export const initS3 = () => {
  const region = process.env.AWS_REGION;
  const accessKeyId = process.env.AWS_ACCESS_KEY_ID;
  const secretAccessKey = process.env.AWS_SECRET_ACCESS_KEY;
  const bucketName = process.env.AWS_S3_BUCKET_NAME;
  const endpoint = process.env.AWS_S3_ENDPOINT; // Optional for MinIO/LocalStack

  console.log(`🔍 S3 Config - Region: ${region}, Bucket: ${bucketName}, Endpoint: ${endpoint || 'default'}`);

  if (!region || !accessKeyId || !secretAccessKey || !bucketName) {
    console.warn('⚠️  Missing AWS S3 configuration. Image upload will not be available.');
    console.warn('   Required: AWS_REGION, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_S3_BUCKET_NAME');
    return;
  }

  s3BucketName = bucketName;

  const config = {
    region,
    credentials: {
      accessKeyId,
      secretAccessKey,
    },
  };

  // Add custom endpoint if specified (for MinIO, LocalStack, etc.)
  if (endpoint) {
    config.endpoint = endpoint;
    config.forcePathStyle = true; // Required for MinIO and some S3-compatible services
    console.log(`🔧 Using custom S3 endpoint: ${endpoint}`);
  }

  s3Client = new S3Client(config);
  console.log('✅ S3 client initialized successfully');
};

/**
 * Get S3 client instance
 * @returns {S3Client|null}
 */
export const getS3Client = () => s3Client;

/**
 * Get S3 bucket name
 * @returns {string|null}
 */
export const getS3BucketName = () => s3BucketName;

export default { initS3, getS3Client, getS3BucketName };
