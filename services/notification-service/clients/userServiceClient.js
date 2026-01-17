/**
 * User Service Client
 * Handles communication with the user-service API
 */

class UserServiceClient {
  constructor() {
    this.baseURL = process.env.USER_SERVICE_URL || 'http://localhost:8082';
  }

  /**
   * Verify JWT token with user-service
   * @param {string} token - JWT token to verify
   * @returns {Promise<Object>} - User information if token is valid
   * @throws {Error} - If token is invalid or request fails
   */
  async verifyToken(token) {
    try {
      const url = `${this.baseURL}/api/user/public/verify`;
      
      const response = await fetch(url, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        const errorData = await response.text();
        throw new Error(`Token verification failed: ${response.status} - ${errorData}`);
      }

      const data = await response.json();
      
      // Extract user info from headers (set by user-service)
      const userId = response.headers.get('X-User-Id');
      const username = response.headers.get('X-Username');
      const userRole = response.headers.get('X-User-Role');

      return {
        success: true,
        userId,
        username,
        userRole,
        data,
      };
    } catch (error) {
      console.error('❌ Token verification error:', error.message);
      throw new Error(`Token verification failed: ${error.message}`);
    }
  }

  /**
   * Get user information by ID
   * @param {string} userId - User ID to fetch
   * @returns {Promise<Object>} - User information
   * @throws {Error} - If request fails
   */
  async getUserById(userId) {
    try {
      const url = `${this.baseURL}/api/user/public/${userId}`;
      
      const response = await fetch(url, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        const errorData = await response.text();
        throw new Error(`Failed to get user: ${response.status} - ${errorData}`);
      }

      const result = await response.json();
      
      // Return the data field from the response
      return result.data || null;
    } catch (error) {
      console.error('❌ Get user error:', error.message);
      throw new Error(`Failed to get user: ${error.message}`);
    }
  }
}

// Export singleton instance
export const userServiceClient = new UserServiceClient();
export default userServiceClient;
