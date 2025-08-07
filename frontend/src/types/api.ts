export interface CreateMatchResponse {
  matchId: string;
  playerId: string;
  token: string;
  invite: string;
}

export interface JoinMatchResponse {
  matchId: string;
  playerId: string;
  token: string;
}
