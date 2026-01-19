/**
 * S3 Upload Service
 * Handle file uploads to AWS S3
 */

import { PutObjectCommand } from '@aws-sdk/client-s3';
import { getS3Client, getS3BucketName } from '../config/s3.js';
import { v4 as uuidv4 } from 'uuid';
import path from 'path';

/**
 * Upload an image to S3
 * @param {Buffer} fileBuffer - File buffer
 * @param {string} originalFilename - Original filename for extension
 * @param {string} mimeType - File MIME type
 * @param {string} folder - S3 folder path (default: 'chat')
 * @returns {Promise<string>} S3 file URL
 */
export const uploadImageToS3 = async (fileBuffer, originalFilename, mimeType, folder = 'chat') => {
  const s3Client = getS3Client();
  const bucketName = getS3BucketName();

  if (!s3Client || !bucketName) {
    throw new Error('S3 client not initialized. Check AWS configuration.');
  }

  // Generate unique filename
  const ext = path.extname(originalFilename);
  const timestamp = Date.now();
  const uniqueId = uuidv4();
  const filename = `${folder}/${uniqueId}-${timestamp}${ext}`;

  const command = new PutObjectCommand({
    Bucket: bucketName,
    Key: filename,
    Body: fileBuffer,
    ContentType: mimeType,
  });

  try {
    await s3Client.send(command);
    
    // Construct the URL
    const region = process.env.AWS_REGION;
    const endpoint = process.env.AWS_S3_ENDPOINT;
    
    // Use custom endpoint if provided, otherwise use standard AWS S3 URL
    let url;
    if (endpoint) {
      // For MinIO or custom S3-compatible services
      url = `${endpoint}/${bucketName}/${filename}`;
    } else {
      // Standard AWS S3 URL
      url = `https://${bucketName}.s3.${region}.amazonaws.com/${filename}`;
    }

    console.log(`✅ Successfully uploaded image to S3: ${filename}`);
    return url;
  } catch (error) {
    console.error('❌ Failed to upload to S3:', error);
    throw new Error(`Failed to upload image to S3: ${error.message}`);
  }
};

export default { uploadImageToS3 };
