import { describe, it, expect } from 'vitest';
import {
  modules,
  getAvailableModules,
  getUpcomingModules,
  getModuleBySlug,
  getTotalLessonCount,
} from './modules';

describe('modules', () => {
  it('should have at least one module', () => {
    expect(modules.length).toBeGreaterThan(0);
  });

  it('every module should have required fields', () => {
    for (const mod of modules) {
      expect(mod.slug).toBeTruthy();
      expect(mod.title).toBeTruthy();
      expect(mod.description).toBeTruthy();
      expect(mod.lessonCount).toBeGreaterThan(0);
      expect(['beginner', 'intermediate', 'advanced']).toContain(mod.difficulty);
    }
  });

  it('should have unique slugs', () => {
    const slugs = modules.map((m) => m.slug);
    expect(new Set(slugs).size).toBe(slugs.length);
  });
});

describe('getAvailableModules', () => {
  it('should exclude upcoming modules', () => {
    const available = getAvailableModules();
    expect(available.every((m) => !m.upcoming)).toBe(true);
  });

  it('should return at least one module', () => {
    expect(getAvailableModules().length).toBeGreaterThan(0);
  });
});

describe('getUpcomingModules', () => {
  it('should only include upcoming modules', () => {
    const upcoming = getUpcomingModules();
    expect(upcoming.every((m) => m.upcoming === true)).toBe(true);
  });
});

describe('getModuleBySlug', () => {
  it('should find an existing module', () => {
    const mod = getModuleBySlug('cli-tools');
    expect(mod).toBeDefined();
    expect(mod!.title).toBe('CLI Tools with Cobra');
  });

  it('should return undefined for non-existent slug', () => {
    expect(getModuleBySlug('does-not-exist')).toBeUndefined();
  });
});

describe('getTotalLessonCount', () => {
  it('should return the sum of all lessons', () => {
    const total = getTotalLessonCount();
    const expected = modules.reduce((sum, m) => sum + m.lessonCount, 0);
    expect(total).toBe(expected);
  });

  it('should be greater than zero', () => {
    expect(getTotalLessonCount()).toBeGreaterThan(0);
  });
});
