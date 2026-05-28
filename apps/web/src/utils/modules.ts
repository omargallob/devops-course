export interface CourseModule {
  slug: string;
  title: string;
  description: string;
  lessonCount: number;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  upcoming?: boolean;
}

export const modules: CourseModule[] = [
  {
    slug: 'cli-tools',
    title: 'CLI Tools with Cobra',
    description: 'Build professional command-line tools using Go and the Cobra library.',
    lessonCount: 6,
    difficulty: 'intermediate',
  },
  {
    slug: 'docker-sdk',
    title: 'Docker SDK for Go',
    description: 'Automate container operations programmatically with the Docker SDK.',
    lessonCount: 5,
    difficulty: 'intermediate',
    upcoming: true,
  },
];

export function getAvailableModules(): CourseModule[] {
  return modules.filter((m) => !m.upcoming);
}

export function getUpcomingModules(): CourseModule[] {
  return modules.filter((m) => m.upcoming);
}

export function getModuleBySlug(slug: string): CourseModule | undefined {
  return modules.find((m) => m.slug === slug);
}

export function getTotalLessonCount(): number {
  return modules.reduce((sum, m) => sum + m.lessonCount, 0);
}
