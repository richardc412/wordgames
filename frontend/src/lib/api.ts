import type { CreateMatchResponse, JoinMatchResponse } from "../types/api";

const API_BASE_URL = "http://localhost:8080";

export const api = {
  async createMatch(): Promise<CreateMatchResponse> {
    const response = await fetch(`${API_BASE_URL}/matches`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`Failed to create match: ${response.statusText}`);
    }

    return response.json();
  },

  async joinMatch(inviteToken: string): Promise<JoinMatchResponse> {
    const response = await fetch(`${API_BASE_URL}/join?invite=${inviteToken}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`Failed to join match: ${response.statusText}`);
    }

    return response.json();
  },
};
