-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY 
    last_fetched_at ASC NULLS FIRST,  -- Сначала фиды, которые никогда не обрабатывались (NULL)
    created_at ASC                    -- Если несколько NULL, берем самый старый
LIMIT 1;