--
-- PostgreSQL database dump
--

-- Dumped from database version 15.8 (Debian 15.8-1.pgdg120+1)
-- Dumped by pg_dump version 15.8

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: films_festivals; Type: TABLE; Schema: public; Owner: xochilpili
--

CREATE TABLE public.films_festivals (
    id integer NOT NULL,
    provider character varying(80) NOT NULL,
    title character varying(250) NOT NULL,
    year smallint DEFAULT 0 NOT NULL,
    genres text[],
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    processed_at timestamp without time zone,
    processed smallint DEFAULT 0
);


ALTER TABLE public.films_festivals OWNER TO xochilpili;

--
-- Name: films_festivals_id_seq; Type: SEQUENCE; Schema: public; Owner: xochilpili
--

CREATE SEQUENCE public.films_festivals_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.films_festivals_id_seq OWNER TO xochilpili;

--
-- Name: films_festivals_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: xochilpili
--

ALTER SEQUENCE public.films_festivals_id_seq OWNED BY public.films_festivals.id;


--
-- Name: films_popular; Type: TABLE; Schema: public; Owner: xochilpili
--

CREATE TABLE public.films_popular (
    id integer NOT NULL,
    provider character varying(80) NOT NULL,
    title character varying(250) NOT NULL,
    year smallint DEFAULT 0 NOT NULL,
    genres text[],
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    processed_at timestamp without time zone,
    processed smallint DEFAULT 0
);


ALTER TABLE public.films_popular OWNER TO xochilpili;

--
-- Name: films_popular_id_seq; Type: SEQUENCE; Schema: public; Owner: xochilpili
--

CREATE SEQUENCE public.films_popular_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.films_popular_id_seq OWNER TO xochilpili;

--
-- Name: films_popular_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: xochilpili
--

ALTER SEQUENCE public.films_popular_id_seq OWNED BY public.films_popular.id;


--
-- Name: tmdb_genres; Type: TABLE; Schema: public; Owner: xochilpili
--

CREATE TABLE public.tmdb_genres (
    id smallint NOT NULL,
    genre character varying(80) NOT NULL
);


ALTER TABLE public.tmdb_genres OWNER TO xochilpili;

--
-- Name: films_festivals id; Type: DEFAULT; Schema: public; Owner: xochilpili
--

ALTER TABLE ONLY public.films_festivals ALTER COLUMN id SET DEFAULT nextval('public.films_festivals_id_seq'::regclass);


--
-- Name: films_popular id; Type: DEFAULT; Schema: public; Owner: xochilpili
--

ALTER TABLE ONLY public.films_popular ALTER COLUMN id SET DEFAULT nextval('public.films_popular_id_seq'::regclass);


--
-- Name: films_festivals films_festivals_pkey; Type: CONSTRAINT; Schema: public; Owner: xochilpili
--

ALTER TABLE ONLY public.films_festivals
    ADD CONSTRAINT films_festivals_pkey PRIMARY KEY (id);


--
-- Name: films_popular films_popular_pkey; Type: CONSTRAINT; Schema: public; Owner: xochilpili
--

ALTER TABLE ONLY public.films_popular
    ADD CONSTRAINT films_popular_pkey PRIMARY KEY (id);


--
-- Name: tmdb_genres tmdb_genres_pkey; Type: CONSTRAINT; Schema: public; Owner: xochilpili
--

ALTER TABLE ONLY public.tmdb_genres
    ADD CONSTRAINT tmdb_genres_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

