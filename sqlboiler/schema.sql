CREATE TABLE pilots (
  id SERIAL NOT NULL PRIMARY KEY,
  name text NOT NULL,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL
);

CREATE TABLE jets (
  id SERIAL NOT NULL PRIMARY KEY,
  pilot_id integer NOT NULL,
  age integer NOT NULL,
  name text NOT NULL,
  color text NOT NULL
);

ALTER TABLE jets ADD CONSTRAINT jet_pilots_fkey FOREIGN KEY (pilot_id) REFERENCES pilots(id) DEFERRABLE;

CREATE TABLE languages (
  id SERIAL NOT NULL PRIMARY KEY,
  language text NOT NULL UNIQUE
);

-- Join table
CREATE TABLE pilots_languages (
  pilot_id integer NOT NULL,
  language_id integer NOT NULL
);

-- Composite primary key
ALTER TABLE pilots_languages ADD CONSTRAINT pilot_language_pkey PRIMARY KEY (pilot_id, language_id);
ALTER TABLE pilots_languages ADD CONSTRAINT pilot_language_pilots_fkey FOREIGN KEY (pilot_id) REFERENCES pilots(id) DEFERRABLE;
ALTER TABLE pilots_languages ADD CONSTRAINT pilot_language_languages_fkey FOREIGN KEY (language_id) REFERENCES languages(id) DEFERRABLE;
