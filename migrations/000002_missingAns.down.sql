ALTER TABLE deployment DROP COLUMN frag_session_missing_completed_at;
ALTER TABLE deployment DROP COLUMN retransmits_completed_at;

ALTER TABLE deployment_device DROP COLUMN all_missing_ans_received;
ALTER TABLE deployment_device DROP COLUMN missing_indices;