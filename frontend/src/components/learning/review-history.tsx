'use client';

import { useState, useEffect } from 'react';
import { api, type ReviewSubmission } from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { History, ChevronDown, ChevronUp, Copy, Check } from 'lucide-react';

interface ReviewHistoryProps {
  userId: string;
}

export const ReviewHistory: React.FC<ReviewHistoryProps> = ({ userId }) => {
  const [reviews, setReviews] = useState<ReviewSubmission[]>([]);
  const [loading, setLoading] = useState(true);
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    loadReviews();
  }, [userId]);

  const loadReviews = async () => {
    try {
      const history = await api.getReviewHistory(userId);
      setReviews(history);
    } catch (err) {
      console.error('Failed to load review history:', err);
    } finally {
      setLoading(false);
    }
  };

  const toggleExpand = (id: string) => {
    setExpandedId(expandedId === id ? null : id);
  };

  const copyCode = async (code: string) => {
    await navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  if (loading) {
    return <div className="text-center py-8">Loading history...</div>;
  }

  if (reviews.length === 0) {
    return (
      <Card>
        <CardContent className="py-8 text-center text-muted-foreground">
          <History className="w-12 h-12 mx-auto mb-4 opacity-50" />
          <p>No reviews yet. Submit code for review to see your history.</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-semibold flex items-center gap-2">
        <History className="w-5 h-5" />
        Review History ({reviews.length})
      </h2>

      {reviews.map((review) => (
        <Card key={review.id} className="overflow-hidden">
          <CardHeader className="py-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Badge variant="outline">{review.exercise_id}</Badge>
                <span className="text-sm text-muted-foreground">
                  {new Date(review.submitted_at).toLocaleDateString()}
                </span>
              </div>
              <Button variant="ghost" size="sm" onClick={() => toggleExpand(review.id)}>
                {expandedId === review.id ? (
                  <ChevronUp className="w-4 h-4" />
                ) : (
                  <ChevronDown className="w-4 h-4" />
                )}
              </Button>
            </div>
          </CardHeader>

          {expandedId === review.id && (
            <CardContent className="pt-0">
              <div className="space-y-4">
                {/* Code Preview */}
                <div className="relative">
                  <pre className="bg-muted p-4 rounded-lg text-sm overflow-x-auto max-h-48">
                    {review.code.slice(0, 500)}
                    {review.code.length > 500 && '...'}
                  </pre>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="absolute top-2 right-2"
                    onClick={() => copyCode(review.code)}
                  >
                    {copied ? <Check className="w-4 h-4" /> : <Copy className="w-4 h-4" />}
                  </Button>
                </div>

                {/* AI Feedback */}
                <div className="bg-purple-50 dark:bg-purple-900/20 border border-purple-200 dark:border-purple-800 rounded-lg p-4">
                  <h4 className="font-semibold text-sm mb-2">AI Feedback</h4>
                  <p className="text-sm whitespace-pre-wrap">{review.feedback}</p>
                </div>
              </div>
            </CardContent>
          )}
        </Card>
      ))}
    </div>
  );
};

export default ReviewHistory;
