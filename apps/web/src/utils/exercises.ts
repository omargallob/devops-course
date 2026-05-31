export interface ExerciseResponse {
  id: string;
  title: string;
  instructions: string;
  starterCode: string;
  hint?: string;
  validationMode: 'exact' | 'regex';
}

export interface ValidateRequest {
  exerciseId: string;
  code: string;
}

export interface ValidateResponse {
  passed: boolean;
  exerciseId: string;
  actualOutput?: string;
  expectedOutput?: string;
  diff?: string;
  compileError?: string;
}

export class ExerciseClient {
  constructor(private baseURL: string = '/api') {}

  async getExercise(exerciseId: string): Promise<ExerciseResponse> {
    const resp = await fetch(`${this.baseURL}/exercises/${exerciseId}`);

    if (!resp.ok) {
      const error = await resp.json().catch(() => ({ error: 'Unknown error' }));
      throw new Error(error.error || `HTTP ${resp.status}`);
    }

    return resp.json();
  }

  async validate(exerciseId: string, code: string): Promise<ValidateResponse> {
    if (!code.trim()) {
      throw new Error('Code cannot be empty');
    }

    const resp = await fetch(`${this.baseURL}/validate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ exerciseId, code } satisfies ValidateRequest),
    });

    if (!resp.ok) {
      const error = await resp.json().catch(() => ({ error: 'Unknown error' }));
      throw new Error(error.error || `HTTP ${resp.status}`);
    }

    return resp.json();
  }
}
