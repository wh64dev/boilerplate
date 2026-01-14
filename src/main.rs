use std::process::exit;

use axum::{Router, routing::get, middleware::from_fn as layer_fn};
use dotenv::dotenv;
use tokio::net::TcpListener;

use crate::{config::Config, controller::index::spa_handler, middleware::log::logger};

mod config;
mod modules;
mod controller;
mod middleware;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    match dotenv() {
        Ok(buf) => buf,
        Err(err) => {
            println!("Error occurred when loading .env files: {:?}", err);
            exit(1);
        }
    };
    
    let config = Config::new();
    println!("Starting {} v{}-{} ({})", config.info.name, config.info.version, config.info.branch, config.info.hash);
    println!("Service build at: {}", config.info.build_time);

    let app = Router::new()
        .fallback(get(spa_handler))
        .layer(layer_fn(logger));

    let server = TcpListener::bind(format!("{}:{}", config.host, config.port))
        .await
        .unwrap();

    println!("webserver binding at: http://{host}:{port}",
        host = config.host,
        port = config.port);
    let _ = axum::serve(server, app).await;
}
