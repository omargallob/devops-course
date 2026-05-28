import { describe, it, expect, vi, beforeEach } from 'vitest';
import { PlaygroundClient } from './playground';

describe('PlaygroundClient', () => {
  let client: PlaygroundClient;

  beforeEach(() => {
    client = new PlaygroundClient('http://localhost:8080/api');
    vi.restoreAllMocks();
  });

  describe('compile', () => {
    it('should throw on empty code', async () => {
      await expect(client.compile('')).rejects.toThrow('Code cannot be empty');
    });

    it('should throw on whitespace-only code', async () => {
      await expect(client.compile('   \n  ')).rejects.toThrow('Code cannot be empty');
    });

    it('should send POST with correct body', async () => {
      const mockResponse = {
        Errors: '',
        Events: [{ Message: 'hello\n', Kind: 'stdout', Delay: 0 }],
        Status: 0,
        IsTest: false,
        TestsFailed: 0,
        VetOK: true,
        VetErrors: '',
      };

      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      });

      const result = await client.compile('package main');

      expect(fetch).toHaveBeenCalledWith('http://localhost:8080/api/compile', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ body: 'package main', version: 2 }),
      });

      expect(result).toEqual(mockResponse);
    });

    it('should throw on HTTP error', async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: false,
        status: 502,
        json: () => Promise.resolve({ error: 'playground unavailable' }),
      });

      await expect(client.compile('package main')).rejects.toThrow('playground unavailable');
    });

    it('should handle non-JSON error response', async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: false,
        status: 500,
        json: () => Promise.reject(new Error('not json')),
      });

      await expect(client.compile('package main')).rejects.toThrow('Unknown error');
    });
  });

  describe('getOutput', () => {
    it('should return errors if present', () => {
      const resp = {
        Errors: 'compile error',
        Events: null,
        Status: 2,
        IsTest: false,
        TestsFailed: 0,
        VetOK: false,
        VetErrors: '',
      };
      expect(client.getOutput(resp)).toBe('compile error');
    });

    it('should concatenate event messages', () => {
      const resp = {
        Errors: '',
        Events: [
          { Message: 'hello ', Kind: 'stdout', Delay: 0 },
          { Message: 'world\n', Kind: 'stdout', Delay: 0 },
        ],
        Status: 0,
        IsTest: false,
        TestsFailed: 0,
        VetOK: true,
        VetErrors: '',
      };
      expect(client.getOutput(resp)).toBe('hello world\n');
    });

    it('should return empty string for null events', () => {
      const resp = {
        Errors: '',
        Events: null,
        Status: 0,
        IsTest: false,
        TestsFailed: 0,
        VetOK: true,
        VetErrors: '',
      };
      expect(client.getOutput(resp)).toBe('');
    });
  });

  describe('isSuccess', () => {
    it('should return true for status 0 with no errors', () => {
      const resp = {
        Errors: '',
        Events: [],
        Status: 0,
        IsTest: false,
        TestsFailed: 0,
        VetOK: true,
        VetErrors: '',
      };
      expect(client.isSuccess(resp)).toBe(true);
    });

    it('should return false for non-zero status', () => {
      const resp = {
        Errors: '',
        Events: [],
        Status: 2,
        IsTest: false,
        TestsFailed: 0,
        VetOK: true,
        VetErrors: '',
      };
      expect(client.isSuccess(resp)).toBe(false);
    });

    it('should return false when errors are present', () => {
      const resp = {
        Errors: 'some error',
        Events: [],
        Status: 0,
        IsTest: false,
        TestsFailed: 0,
        VetOK: true,
        VetErrors: '',
      };
      expect(client.isSuccess(resp)).toBe(false);
    });
  });
});
