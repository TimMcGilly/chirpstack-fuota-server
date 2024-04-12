ALTER TABLE deployment ADD COLUMN frag_session_missing_completed_at  timestamp with time zone null;
ALTER TABLE deployment ADD COLUMN retransmits_completed_at timestamp with time zone null;

ALTER TABLE deployment_device ADD COLUMN all_missing_ans_received timestamp with time zone null;
ALTER TABLE deployment_device ADD COLUMN missing_indices text;