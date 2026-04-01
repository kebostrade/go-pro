export default function Loading() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white dark:from-gray-900 dark:to-gray-800">
      <div className="container mx-auto px-4 py-8">
        <div className="animate-pulse space-y-8">
          {/* Hero skeleton */}
          <div className="h-48 bg-gray-200 dark:bg-gray-700 rounded-2xl" />
          {/* Tabs skeleton */}
          <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg overflow-hidden">
            <div className="flex border-b border-gray-200 dark:border-gray-700">
              <div className="flex-1 h-14 bg-gray-100 dark:bg-gray-700" />
              <div className="flex-1 h-14 bg-gray-50 dark:bg-gray-750" />
              <div className="flex-1 h-14 bg-gray-50 dark:bg-gray-750" />
            </div>
            <div className="p-8 space-y-6">
              {/* Overview tab skeleton */}
              <div className="space-y-4">
                <div className="h-8 w-48 bg-gray-200 dark:bg-gray-700 rounded" />
                <div className="flex gap-2">
                  <div className="h-8 w-24 bg-gray-200 dark:bg-gray-700 rounded-lg" />
                  <div className="h-8 w-32 bg-gray-200 dark:bg-gray-700 rounded-lg" />
                  <div className="h-8 w-28 bg-gray-200 dark:bg-gray-700 rounded-lg" />
                </div>
              </div>
              <div className="space-y-3">
                <div className="h-6 w-40 bg-gray-200 dark:bg-gray-700 rounded" />
                <div className="h-20 bg-gray-100 dark:bg-gray-750 rounded-lg" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
