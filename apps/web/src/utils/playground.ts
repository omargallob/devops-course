export interface CompileRequest {
  body: string;
  version?: number;
}

export interface CompileEvent {
  Message: string;
  Kind: string;
  Delay: number;
}

export interface CompileResponse {
  Errors: string;
  Events: CompileEvent[] | null;
  Status: number;
  IsTest: boolean;
  TestsFailed: number;
  VetOK: boolean;
  VetErrors: string;
}

export class PlaygroundClient {
  constructor(private baseURL: string = '/api') {}

  async compile(code: string): Promise<CompileResponse> {
    if (!code.trim()) {
      throw new Error('Code cannot be empty');
    }

    const resp = await fetch(`${this.baseURL}/compile`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ body: code, version: 2 } satisfies CompileRequest),
    });

    if (!resp.ok) {
      const error = await resp.json().catch(() => ({ error: 'Unknown error' }));
      throw new Error(error.error || `HTTP ${resp.status}`);
    }

    return resp.json();
  }

  getOutput(response: CompileResponse): string {
    if (response.Errors) {
      return response.Errors;
    }
    if (!response.Events) {
      return '';
    }
    return response.Events.map((e) => e.Message).join('');
  }

  isSuccess(response: CompileResponse): boolean {
    return response.Status === 0 && !response.Errors;
  }
}
