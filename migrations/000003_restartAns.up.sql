ALTER TABLE deployment ADD COLUMN restart_completed_at  timestamp with time zone null;

ALTER TABLE deployment_device ADD COLUMN restart_completed_at timestamp with time zone null;
