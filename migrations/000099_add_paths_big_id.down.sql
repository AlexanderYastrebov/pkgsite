-- Copyright 2021 The Go Authors. All rights reserved.
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

BEGIN;

ALTER TABLE paths DROP COLUMN big_id;
DROP TRIGGER set_paths_big_id ON paths;
DROP FUNCTION set_paths_big_id;

END;
