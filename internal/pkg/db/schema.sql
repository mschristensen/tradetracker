--
-- PostgreSQL database dump
--

-- Dumped from database version 12.10
-- Dumped by pg_dump version 12.10

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
-- Name: migrations; Type: TABLE; Schema: public; Owner: tradetracker
--

CREATE TABLE public.migrations (
    id text NOT NULL,
    applied_at timestamp with time zone
);


ALTER TABLE public.migrations OWNER TO tradetracker;

--
-- Name: positions; Type: TABLE; Schema: public; Owner: tradetracker
--

CREATE TABLE public.positions (
    id integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    instrument_id bigint NOT NULL,
    size bigint NOT NULL,
    "timestamp" timestamp without time zone NOT NULL
);


ALTER TABLE public.positions OWNER TO tradetracker;

--
-- Name: positions_id_seq; Type: SEQUENCE; Schema: public; Owner: tradetracker
--

CREATE SEQUENCE public.positions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.positions_id_seq OWNER TO tradetracker;

--
-- Name: positions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: tradetracker
--

ALTER SEQUENCE public.positions_id_seq OWNED BY public.positions.id;


--
-- Name: trades; Type: TABLE; Schema: public; Owner: tradetracker
--

CREATE TABLE public.trades (
    id integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    instrument_id bigint NOT NULL,
    size bigint NOT NULL,
    price numeric NOT NULL,
    "timestamp" timestamp without time zone NOT NULL
);


ALTER TABLE public.trades OWNER TO tradetracker;

--
-- Name: trades_id_seq; Type: SEQUENCE; Schema: public; Owner: tradetracker
--

CREATE SEQUENCE public.trades_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.trades_id_seq OWNER TO tradetracker;

--
-- Name: trades_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: tradetracker
--

ALTER SEQUENCE public.trades_id_seq OWNED BY public.trades.id;


--
-- Name: positions id; Type: DEFAULT; Schema: public; Owner: tradetracker
--

ALTER TABLE ONLY public.positions ALTER COLUMN id SET DEFAULT nextval('public.positions_id_seq'::regclass);


--
-- Name: trades id; Type: DEFAULT; Schema: public; Owner: tradetracker
--

ALTER TABLE ONLY public.trades ALTER COLUMN id SET DEFAULT nextval('public.trades_id_seq'::regclass);


--
-- Name: migrations migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: tradetracker
--

ALTER TABLE ONLY public.migrations
    ADD CONSTRAINT migrations_pkey PRIMARY KEY (id);


--
-- Name: positions positions_pkey; Type: CONSTRAINT; Schema: public; Owner: tradetracker
--

ALTER TABLE ONLY public.positions
    ADD CONSTRAINT positions_pkey PRIMARY KEY (id);


--
-- Name: trades trades_pkey; Type: CONSTRAINT; Schema: public; Owner: tradetracker
--

ALTER TABLE ONLY public.trades
    ADD CONSTRAINT trades_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

